#!/bin/bash

echo "Testing sunaba_code sandbox..."

# Test 1: Try to write to current directory
echo "Test 1: Writing to current directory (should succeed)"
./sunaba_code -- /bin/bash -c "echo 'test' > test_file.txt && cat test_file.txt && rm test_file.txt"

echo -e "\nTest 2: Try to write to /tmp (should fail)"
./sunaba_code -- /bin/bash -c "echo 'test' > /tmp/test_file.txt" || echo "Failed as expected"

echo -e "\nTest 3: Try to write to /tmp with -w flag (should succeed)"
./sunaba_code -w /tmp -- /bin/bash -c "echo 'test' > /tmp/test_file.txt && cat /tmp/test_file.txt && rm /tmp/test_file.txt"

echo -e "\nTest 4: Network test without --network flag (should fail)"
./sunaba_code -- /usr/bin/curl -s https://api.github.com || echo "Network access blocked as expected"

echo -e "\nTest 5: Network test with --network flag (should succeed)"
./sunaba_code --network -- /usr/bin/curl -s https://api.github.com | head -5

echo -e "\nDone!"