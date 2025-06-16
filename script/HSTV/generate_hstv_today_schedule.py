#!/usr/bin/env python3

import csv
import datetime
from pathlib import Path


def update_dates(csv_file, output_file=None):
    # Get today's date and tomorrow's date in the format MM-DD-YYYY
    today_dt = datetime.datetime.now()
    today = today_dt.strftime("%m-%d-%Y")
    tomorrow_dt = today_dt + datetime.timedelta(days=1)
    tomorrow = tomorrow_dt.strftime("%m-%d-%Y")

    # If no output file is specified, create one with today's date
    if output_file is None:
        output_file = csv_file.parent / f"schedule_{today}.csv"

    # Read the CSV file
    rows = []
    with open(csv_file, "r") as file:
        reader = csv.reader(file)
        header = next(reader)  # Get the header
        rows.append(header)

        # Process each row
        for row in reader:
            if len(row) >= 5:  # Ensure the row has enough columns
                # Update Start Date (index 1) to today
                row[1] = today

                # Special case: if End Time is 00:00:00, set End Date to tomorrow
                # Otherwise, set End Date to today
                if row[4] == "00:00:00":
                    row[3] = tomorrow
                else:
                    row[3] = today
            rows.append(row)

    # Write the updated data to the output file
    with open(output_file, "w", newline="") as file:
        writer = csv.writer(file)
        writer.writerows(rows)

    print(f"Created new file {output_file} with updated dates ({today})")
    return output_file


if __name__ == "__main__":
    # Path to the CSV file - use the script's directory to find the CSV
    script_dir = Path(__file__).parent
    csv_file = script_dir / "schedule.csv"

    # Update the dates and create a new file with today's date in the filename
    new_file = update_dates(csv_file)
    print(f"New file created: {new_file}")
