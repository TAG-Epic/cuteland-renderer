from random import choice
from json import dumps

d = []

for x in range(5):
    d.append([])
    for y in range(5):
        d[x].append(choice([None, None, "tile:ore_fresh", "tile:ore_used"]))

print(dumps(d))
