Simple monitoring application that accepts a url and sends a specified amount of test traffic to it calculating the number of failed requests, useful for finding intermittant failures
While it may be odd that it prompts for tcp or udp on a curl my intention is to have this usable for tcp/udp endpoint checking over a variety of ports


run with the compiled go or build after mofidying 

command line arguement is 
./curltr https://websiteyouaretesting.com

answer the prompts
