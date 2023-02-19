# Population
To create or to validate a popultation a certain amount of conditions have to be met.  
For example the cut has to be a guillotine cut. No item shall expand to borders of the stock plate.  
The total surface area of the items cannot be bigger thant the stock plates etc.  

Ideas:
- if there is just one cut trhough the length of the stock plate possible, we have a guillocut

# Fitness Function
The three criterias are:
- minimize waste (c1)
- minimize number of stock plates (c2)
- minimize stages (c3)

## Minimizing waste
For minimizing the waste the general idea of the fitness function would be:  
- As: Total area of all stock plates
- Ai: Total area of all items
- w: waste  

*As - (Ai + w) = 100 for w = 0*
*As - (Ai + w) = 0 for w = As - Ai*

If the waste is 0 we do already use the least amount of stock plates.  
So i think that the number of used stock plates does not have to be considered in the fitness function.

## Minimzing number of stock plates
The goal should always be to just use 1 stock plate.  
So 1 stock plate = 100% fitness  
Since the number of stock plates is variable the following function is needed.
- Sn number of stock plates used  
- St total number of stock plates  
*(St / Sn) * (100 / St)*  

## Minimizing stages (n)
For the number of stages (n) we assume a maximum of 20 stages.  
1 stage can only appear is we have a maximum of 2 items this case can be ignored.  
Even 2 stages are highly unlikely in everyday situations.  
We assume that 3 stages and below are a perfect outcome.  
So 3 stages are 100% fitness and everything above 20 stages is 0% fitness.

## Combine waste and stage
To combine the 2 criteria we need to give those a weight in percent. For example:
- c1 = 45% = 0.45
- c2 = 45% = 0.45
- c3 = 10% = 0.1

This implies that a perfect solution would be a waste of 0 and less or equal than 3 stages.  
At this point we have 100% overall fitness and can succesfully end the loop.  

Total fitness function could look like this:  
`((As - (Ai + w)) * c1) + ((St / Sn) * (100 / St) * c2) + (n * c3) = 1; for w = 0, Sn = 1 and n <= 3`
`((As - (Ai + w)) * c1) + ((St / Sn) * (100 / St) * c2) + (n * c3) = 0; for w = As - Ai and n >= 20`
with:
- As: total surface area of stock plates
- Ai: total surface area of items
- w: surface area of waste
- St: total number of stock plates
- Sn: number of used stock plates
- n: number of stages
- c1: wighting of minimal waste in percent
- c2: wighting of minimal number of stock plates in percent
- c3: wighting of minimal number of stages in percent

# (Natural) Selection
do not just use the 2 best. make it based on probalities.

# Crossover

# Mutation

# TODO
- think about natural selection
- think about how to implement crossover
- think about hot wo implement mutation
- think about exponential fitness function
