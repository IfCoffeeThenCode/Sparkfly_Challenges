# Sparkfly_Challenges
Code challenge for application to SparkFly developer position

## Requirements

1. We have generated a number of (csv) code files that should be guaranteed to have unique codes (e.g. X59J) not only within themselves, but also among each other. Write code that can take the names of these files and verify the code uniqueness across all csv files in parallel. Bonus points if you can make it immediately stop itself once it has found a duplicate code. Example csv files can be found in the attachments of this email. This solution does NOT need to have tests and should be at most a few hundred lines. Please email me if your solution grows beyond this limited scope.

2. We receive large (20MB+) code files that have to be stored in S3 for record-keeping. To minimize costs, we would like to store them in a compressed format. To further minimize costs, we would like to offload this process onto low-memory hardware. We get these files regularly and need the software that processes them to be expedient. For simplicity, we have decided to use the gzip compression format as it offers the balance between speed/compression that we need. Please write code that takes uncompressed input and writes compressed output and test(s) that verify its efficacy. The interface requirements are:

a. The upload manager to S3 takes an io.Reader as its argument (output from your code)
b. The uncompressed data is provided to your code as an io.ReadCloser (input to your code)
You are encouraged to mock out these inputs and outputs to simplify your solution
