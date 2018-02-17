# xdcc2

go backend for searching anime and manga packlists

temporary hosted server here https://immense-beyond-13018.herokuapp.com/

2/17/2018 Update
I actually did some algorithm studies and I realized that I used some pretty inefficient data structures in some areas of this project. Most notably, the big perpetrator that irritates me slightly is a triple O(n^3) in the optimize downloads file. While it is an inefficient runtime, the actual slice that is mapped over is only proportional to slice chunk size. Currently, I've set that to two. To be honest, most of the data searches I've utilized iterate over exceptionally small lengths (length of 1 or 2),

However, it's not very scalable if I do, say, a chunk of 1000, or 10,000. While it's in the realm of theoretical discussion, it's someone that I need to keep in mind going forward if I want to elevate my efficiency. Due to the current performance which I'm satisfied with, I'll leave it as is. I could change it to O(1) with a map, but it's only creating a map of one key-value pair. I should use maps more next time, though.

There are some more front-end changes coming, such as animation styling and directly queuing straight from the homepage. There might also be a small cache error in the front end of electron.
