[link](https://fly.io/dist-sys/2/)

This one I solved using google library uuid but this could be solved by having 
a random string generator. 

Why not a number? numbers are not unique, and if you have a lot of nodes,
you will have a lot o collision and you will have to deal with it.
Prefer a string, it's unique and you don't hit the limit of the number of
nodes you can have.