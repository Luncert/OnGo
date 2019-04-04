
import random

num = 10000000
a, b = 1000, 1000000
with open('data', 'wb') as f:
    for _ in range(num - 1):
        f.write(str(random.randint(a, b)).encode('utf-8'))
        f.write(','.encode('utf-8'))
    f.write(str(random.randint(a, b)).encode('utf-8'))
    
