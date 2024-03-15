# Changes Summary

-	i)		Journeys is now a map indexed by Id


-	ii)		Cars is now a map indexed by Id


-	iii)	A new array of car slices indexed by the number of available seats (CarsByAvailableSeats) has been added to the Carpool service.


-	iv)		When assigning a new journey to a car, the carpool service will firstly try to assign it to a car having the exact number of available seats.
     if it's not possible, it will apply an allocation strategy based on the worst fit, trying to leave the largest number of seats available in the car (which could be used for subsequent journeys).


-	v)		The Carpool Reassign function won't stop after finding the first pending journey that would fit into the car(as it was doing in the original version of the function),
     the function will now loop through all the pending journeys while there's still available seats in the car.
 
## A few words about the changes

- 	CarPool's Pending has been kept has a slice of journey pointers: Since journeys arrival order has to be taken into account (when reassigning), the slice has to be traversed from the first element to the last which seems fit to this structure.
    However, if a PostDropoff request is made for a car stored in pending, a sequential search will be performed in order to find the journey.
    Also, deleting it from the list can be an expensive operation as it will require in most cases shifting (copying) part of the slice content.
    It has to be tested with a bigger amount of data but replacing it with a linked list should result in reducing the cost of the deleting operation.


- 	About the Worst Strategy algorithm: I think it is an algorithm that could work nicely since it tries to leave the biggest amount of available seats (that could be used for other journeys). 
    However, it is not necessarily the best strategy (for instance, it may make it difficult for a 6 people journey to find a car). Still, it has to be tested with real data.
    Still, it has to be tested with real data to actually see how good it would really work. I've also included a second allocation strategy (BestFit).


- 	No id index has been included in CarsByAvailableSeats. Having every item in this array indexed by Id would prevent the Carpool service to perform a sequential search while deleting a *car during dropoff y reassign.
    If, when tested, these delete operations prove too slow, the array could be updated by indexing by id each of the array elements.


- 	CarPool's Cars has been indexed by id foreseeing more API requests involving cars. 
    The presence of an id in the car model makes me think that it has - or will have - greater importance in the carpooling system and, since I've anticipated that more request (having car's id as a parameter) will come, I have decided to index them.

    
- About change #v, the changes made to the reassign function, a minGroup or an array including a counter of pending journeys per person could be used as an early out condition of the loop when there's no pending journey that would fit into the car.

# Some thoughts about the Model:

- 	The more I worked on the coding challenge, the more the idea of a journey only as a group of people already assigned to a car came to mind.
    I see the group of people with the intention of traveling as an entity different from the journey (I think the name Party could fit to this new entity). 
    In this case, Journey would simply be a struct with a first pointer to a car and second to this new Party entity.	


- 	Having a Pending list of Party elements seems to make more conceptual sense than a list of pending journeys. I think this Party entity would have more information associated with it in a more complex system (for example, a pickup time and location).
    However, taking into account the simplified scenario that has been provided, I have decided to essentially keep the proposed models (except for some String functions and the necessary corrections to pass the acceptance tests).