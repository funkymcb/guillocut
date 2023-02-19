# Population
To create or to validate a popultation a certain amount of conditions have to be met.  
For example the cut has to be a guillotine cut. No item shall expand to borders of the stock plate.  
The total surface area of the items cannot be bigger thant the stock plates etc.

# Fitness Function
The three criterias are:
    - c1: minimize waste
    - c2: minimize stages

## Minimizing waste
For minimizing the waste the general idea of the fitness function would be:  
- As: Total area of all stock plates
- Ai: Total area of all items
- w: waste
*As - Ai ^= As - (Ai + w)*  
- If *w = 0* the fitness is 100%
- If *w = As - Ai* the fitness is 0%

If the waste is 0 we do already use the least amount of stock plates.  
So i think that the number of used stock plates does not have to be considered in the fitness function.

## Minimizing stages
For the number of stages we assume a maximum of 20 stages.  
1 stage can only appear is we have a maximum of 2 items this case can be ignored.  
Even 2 stages are highly unlikely in everyday situations.  
We assume that 3 stages and below are a perfect outcome.  
So 3 stages are 100% fitness and everything above 20 stages is 0% fitness.

## Combine waste and stage
To combine the 3 criteria we use the following percentage:
- c1 = 70%
- c2 = 30%

This implies that a perfect solution would be a waste of 0 and less or equal than 3 stages.  
At this point we have 100% overall fitness and can succesfully end the loop.

# TODO
- think about number of stock plates again
