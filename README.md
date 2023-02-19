# guillocut
GuilloCut is a **blazingly fast** tool for solving the two-dimensional multi-stage cutting stock problem

## Problem description
As mentioned previous the main algorithm solves the two-dimensional multi-stage cutting stock problem.  
It's goal is to cut n number of rectangular items (i) from m number of rectangular stock plates (I).  
The stock plate usage (waste) shall be minimized and the cuts shall be limited to guillotine cuts.  
There are no restrictions on how many stages (90° rotations) these guillotine cuts can be performed.  
However the number of stages shall also be minimized aswell.

## Approach
The algorithm follows a genetic approach. Check approach.md for more details.

## ToDos
Implement:
- [ ] canvas for drawing population
- [ ] pupulation creation
    - [ ] general validation algorithm
    - [ ] guillotine cut validation algorithm
- [ ] fitness function
    - [ ] minimum waste algorithm
    - [ ] stage algorithm
    - [ ] combination of the two
