i=0
n=10000
while [[ $i -le $n ]]
do   
    a=$RANDOM
    let "i++"
    echo "{\"ID\":\"hello$a\"}"
    curl localhost:8080/playback -d "{\"ID\":\"hello$a\"}"
done 
