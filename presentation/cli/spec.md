# CLI for memory usage.

Provide memory monitor for golang map building.

Make a `[]map[int]any`.

Cli arguments for:

* number for the number of maps created with (make).
* number for populating the number of key value pairs of the map.
* Each map element will be made up of a key (int) and struct{}{}.

Aim:

* We want to see what the allocation of memory for the program is:
    * before any maps are initialised.
    * once the maps have been initialised.
    * once the maps are populated with elements.


    

