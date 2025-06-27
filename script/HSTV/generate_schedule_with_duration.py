#!/usr/bin/env python3

import csv
import datetime
import random
import os
from pathlib import Path
import sys


def generate_schedule_with_duration(input_file=None, output_file=None, date_str=None):
    """
    Generate a schedule CSV file using content IDs and durations from contentids_duration.csv.

    Args:
        input_file: Path to the input file containing content IDs and durations.
        output_file: Path to the output file. If None, a default name will be used.
        date_str: Date string in MM-DD-YYYY format. If None, today's date will be used.

    Returns:
        Path to the generated file.
    """
    # Set default input file if not provided
    if input_file is None:
        script_dir = Path(__file__).parent
        input_file = script_dir / "contentids_duration.csv"

    # Read content IDs and durations from the input file
    content_data = []
    with open(input_file, "r") as file:
        reader = csv.reader(file)
        next(reader)  # Skip header
        for row in reader:
            if len(row) >= 2:
                content_id = row[0]
                duration = int(row[1])
                content_data.append((content_id, duration))

    # Get the date for the schedule
    if date_str is None:
        today_dt = datetime.datetime.now()
        date_str = today_dt.strftime("%m-%d-%Y")
    else:
        # Parse the provided date string to ensure it's valid
        try:
            today_dt = datetime.datetime.strptime(date_str, "%m-%d-%Y")
        except ValueError:
            print(f"Invalid date format: {date_str}. Using today's date instead.")
            today_dt = datetime.datetime.now()
            date_str = today_dt.strftime("%m-%d-%Y")

    # Calculate tomorrow's date
    tomorrow_dt = today_dt + datetime.timedelta(days=1)
    tomorrow_str = tomorrow_dt.strftime("%m-%d-%Y")

    # If no output file is specified, create one with the date
    if output_file is None:
        script_dir = Path(__file__).parent
        output_file = script_dir / f"schedule_{date_str.replace('-', '_')}.csv"

    # Initialize the schedule with header
    schedule = [["ContentId", "Start Date", "Start Time", "End Date", "End Time"]]

    # Start time is midnight (00:00:00)
    current_time = datetime.datetime.combine(today_dt.date(), datetime.time(0, 0, 0))
    end_of_day = datetime.datetime.combine(tomorrow_dt.date(), datetime.time(0, 0, 0))

    # For 0-9 o'clock, keep only one row
    if current_time.hour < 9:
        # Randomly select a content ID and its duration
        content_id, duration = random.choice(content_data)

        # Calculate end time (9:00:00)
        entry_end_time = datetime.datetime.combine(
            today_dt.date(), datetime.time(9, 0, 0)
        )

        # Format the times for the CSV
        start_date = current_time.strftime("%m-%d-%Y")
        start_time = current_time.strftime("%H:%M:%S")
        end_date = entry_end_time.strftime("%m-%d-%Y")
        end_time_str = entry_end_time.strftime("%H:%M:%S")

        # Add the entry to the schedule
        schedule.append([content_id, start_date, start_time, end_date, end_time_str])

        # Update current time
        current_time = entry_end_time

    # Generate schedule entries until we reach the end of the day
    while current_time < end_of_day:
        # Calculate remaining time until end of day
        remaining_seconds = (end_of_day - current_time).total_seconds()

        # If there's less than 5 minutes left, extend the previous entry to midnight
        if remaining_seconds < 300 and len(schedule) > 1:
            # Update the end time of the last entry to midnight
            schedule[-1][3] = (
                tomorrow_str if end_of_day.date() > current_time.date() else start_date
            )
            schedule[-1][4] = end_of_day.strftime("%H:%M:%S")
            break

        # Find content that fits in the remaining time
        suitable_content = [
            (cid, dur) for cid, dur in content_data if dur <= remaining_seconds
        ]

        # If no suitable content found, extend the previous entry to midnight
        if not suitable_content and len(schedule) > 1:
            # Update the end time of the last entry to midnight
            schedule[-1][3] = (
                tomorrow_str if end_of_day.date() > current_time.date() else start_date
            )
            schedule[-1][4] = end_of_day.strftime("%H:%M:%S")
            break

        # Randomly select a content ID and its duration that fits
        content_id, duration = random.choice(
            suitable_content if suitable_content else content_data
        )

        # Calculate a random additional time between 5-10 minutes (300-600 seconds)
        additional_time = random.randint(300, 900)

        # Total duration in seconds (ensure it's at least the content duration)
        total_duration = max(duration, additional_time)

        # Make sure the total duration doesn't exceed the remaining time
        if total_duration > remaining_seconds:
            total_duration = remaining_seconds

        # Convert duration to timedelta
        duration_delta = datetime.timedelta(seconds=total_duration)

        # Calculate the end time for this entry
        entry_end_time = current_time + duration_delta

        # Format the times for the CSV
        start_date = current_time.strftime("%m-%d-%Y")
        start_time = current_time.strftime("%H:%M:%S")

        # Determine if the end time is on the next day
        if entry_end_time.date() > current_time.date():
            end_date = tomorrow_str
        else:
            end_date = start_date

        end_time_str = entry_end_time.strftime("%H:%M:%S")

        # Add the entry to the schedule
        schedule.append([content_id, start_date, start_time, end_date, end_time_str])

        # Move to the next time slot
        current_time = entry_end_time

    # Write the schedule to the output file
    with open(output_file, "w", newline="") as file:
        writer = csv.writer(file)
        writer.writerows(schedule)

    print(f"Created schedule file: {output_file}")
    return output_file


def generate_schedules_for_days(num_days=7, start_date=None):
    """
    Generate schedules for multiple days.

    Args:
        num_days: Number of days to generate schedules for.
        start_date: Start date in MM-DD-YYYY format. If None, today's date will be used.

    Returns:
        List of paths to the generated files.
    """
    # Get the start date
    if start_date is None:
        start_dt = datetime.datetime.now()
    else:
        try:
            start_dt = datetime.datetime.strptime(start_date, "%m-%d-%Y")
        except ValueError:
            print(f"Invalid date format: {start_date}. Using today's date instead.")
            start_dt = datetime.datetime.now()

    # Set paths
    script_dir = Path(__file__).parent
    input_file = script_dir / "duration_content_id.csv"

    # Generate schedules for each day
    output_files = []
    for i in range(num_days):
        # Calculate the date for this schedule
        current_dt = start_dt + datetime.timedelta(days=i)
        current_date = current_dt.strftime("%m-%d-%Y")

        # Set output file path
        output_file = (
            script_dir / f"schedule_random_{current_date.replace('-', '_')}.csv"
        )

        # Generate schedule for this day
        generated_file = generate_schedule_with_duration(
            input_file=input_file, output_file=output_file, date_str=current_date
        )

        output_files.append(generated_file)

    return output_files


def generate_combined_schedule(num_days=7, start_date=None, output_file=None):
    """
    Generate a single CSV file containing schedules for multiple days.

    Args:
        num_days: Number of days to generate schedules for.
        start_date: Start date in MM-DD-YYYY format. If None, today's date will be used.
        output_file: Path to the output file. If None, a default name will be used.

    Returns:
        Path to the generated file.
    """
    # Get the start date
    if start_date is None:
        start_dt = datetime.datetime.now()
    else:
        try:
            start_dt = datetime.datetime.strptime(start_date, "%m-%d-%Y")
        except ValueError:
            print(f"Invalid date format: {start_date}. Using today's date instead.")
            start_dt = datetime.datetime.now()

    # Set paths
    script_dir = Path(__file__).parent
    input_file = script_dir / "duration_content_id.csv"

    # If no output file is specified, create one with the date range
    if output_file is None:
        end_dt = start_dt + datetime.timedelta(days=num_days - 1)
        start_str = start_dt.strftime("%m_%d_%Y")
        end_str = end_dt.strftime("%m_%d_%Y")
        output_file = script_dir / f"schedule_combined_{start_str}_to_{end_str}.csv"

    # Read content IDs and durations from the input file
    content_data = []
    with open(input_file, "r") as file:
        reader = csv.reader(file)
        next(reader)  # Skip header
        for row in reader:
            if len(row) >= 2:
                content_id = row[0]
                duration = int(row[1])
                content_data.append((content_id, duration))

    # Initialize the combined schedule with header
    combined_schedule = [
        ["ContentId", "Start Date", "Start Time", "End Date", "End Time"]
    ]

    # Generate schedule for each day and add to the combined schedule
    for i in range(num_days):
        # Calculate the date for this schedule
        current_dt = start_dt + datetime.timedelta(days=i)
        current_date = current_dt.strftime("%m-%d-%Y")

        # Calculate tomorrow's date
        tomorrow_dt = current_dt + datetime.timedelta(days=1)
        tomorrow_str = tomorrow_dt.strftime("%m-%d-%Y")

        # Start time is midnight (00:00:00)
        current_time = datetime.datetime.combine(
            current_dt.date(), datetime.time(0, 0, 0)
        )
        end_of_day = datetime.datetime.combine(
            tomorrow_dt.date(), datetime.time(0, 0, 0)
        )

        # For 0-9 o'clock, keep only one row
        if current_time.hour < 9:
            # Randomly select a content ID and its duration
            content_id, duration = random.choice(content_data)

            # Calculate end time (9:00:00)
            entry_end_time = datetime.datetime.combine(
                current_dt.date(), datetime.time(9, 0, 0)
            )

            # Format the times for the CSV
            start_date = current_time.strftime("%m-%d-%Y")
            start_time = current_time.strftime("%H:%M:%S")
            end_date = entry_end_time.strftime("%m-%d-%Y")
            end_time_str = entry_end_time.strftime("%H:%M:%S")

            # Add the entry to the schedule
            combined_schedule.append(
                [content_id, start_date, start_time, end_date, end_time_str]
            )

            # Update current time
            current_time = entry_end_time

        # Generate schedule entries until we reach the end of the day
        while current_time < end_of_day:
            # Calculate remaining time until end of day
            remaining_seconds = (end_of_day - current_time).total_seconds()

            # If there's less than 5 minutes left, extend the previous entry to midnight
            if remaining_seconds < 300 and len(combined_schedule) > 1:
                # Update the end time of the last entry to midnight
                combined_schedule[-1][3] = tomorrow_str
                combined_schedule[-1][4] = end_of_day.strftime("%H:%M:%S")
                break

            # Find content that fits in the remaining time
            suitable_content = [
                (cid, dur) for cid, dur in content_data if dur <= remaining_seconds
            ]

            # If no suitable content found, extend the previous entry to midnight
            if not suitable_content and len(combined_schedule) > 1:
                # Update the end time of the last entry to midnight
                combined_schedule[-1][3] = tomorrow_str
                combined_schedule[-1][4] = end_of_day.strftime("%H:%M:%S")
                break

            # Randomly select a content ID and its duration that fits
            content_id, duration = random.choice(
                suitable_content if suitable_content else content_data
            )

            # Calculate a random additional time between 5-10 minutes (300-600 seconds)
            additional_time = random.randint(300, 900)

            # Total duration in seconds (ensure it's at least the content duration)
            total_duration = max(duration, additional_time)

            # Make sure the total duration doesn't exceed the remaining time
            if total_duration > remaining_seconds:
                total_duration = remaining_seconds

            # Convert duration to timedelta
            duration_delta = datetime.timedelta(seconds=total_duration)

            # Calculate the end time for this entry
            entry_end_time = current_time + duration_delta

            # Format the times for the CSV
            start_date = current_time.strftime("%m-%d-%Y")
            start_time = current_time.strftime("%H:%M:%S")

            # Determine if the end time is on the next day
            if entry_end_time.date() > current_time.date():
                end_date = tomorrow_str
            else:
                end_date = start_date

            end_time_str = entry_end_time.strftime("%H:%M:%S")

            # Add the entry to the schedule
            combined_schedule.append(
                [content_id, start_date, start_time, end_date, end_time_str]
            )

            # Move to the next time slot
            current_time = entry_end_time

    # Write the combined schedule to the output file
    with open(output_file, "w", newline="") as file:
        writer = csv.writer(file)
        writer.writerows(combined_schedule)

    print(f"Created combined schedule file: {output_file}")
    return output_file


def main():
    """
    Main function to generate a combined schedule for the next 7 days.
    """
    # Get today's date
    today = datetime.datetime.now()
    today_str = today.strftime("%m-%d-%Y")
    # today_str = "06-30-2025"
    if len(sys.argv) > 1:
        today_str = sys.argv[1]
    num_days = 1

    # Generate a combined schedule for the next 7 days
    output_file = generate_combined_schedule(num_days=num_days, start_date=today_str)

    print(
        f"Generated combined schedule for {num_days} days starting from {today_str} at {output_file}"
    )


if __name__ == "__main__":
    main()
