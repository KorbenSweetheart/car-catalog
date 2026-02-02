# viewer

A website that showcases information about different car models, their specifications, manufacturers, and more.

Mandatory:
[x] Catalog page
[x] Advanced filtering with different search options.
[x] Display manufacturer info on a car page
[x] Add car search
[x] Add a proper 404 errors handling (DRY)
[x] Comparison of different car models in terms of features and specifications.
[ ] Implement recommendations based on the visited cars
[x] Add maintainance page
[ ] Implements other features that aren't listed in the bonus part.

TODO:
[x] Add Popular body types for homepage
[x] Replace the icons on popular body types
[x] implement a hamburger menu for mobile devices
[x] Add favicon
[ ] Fix mobile view of the compare page
[ ] add reserve a car page with form
[x] limit search input string to 50
[ ] Update filter/catalog handler to remove empty filter parameters from query
[ ] make friendly query parameters in compare hlml page
[ ] sort manufacturers in filters
[ ] maybe use pointers in cache for performance.
[ ] update for loops to use index instead of copy the object
[ ] maybe we have to swap Transmission and Gearbox, so transmission will display native data, and gearbox would be used to filter, and not displayed on a page
[ ] refactor CSS

Maybe:
[ ] Update carapi images, prepare a better one with good resolution and ratio
[ ] Maybe each usecase file should be a separate usecase struct.
[ ] Put everyting in Docker
[ ] Try Redis as a repo or cache replacement
[ ] Try SQLite as a repo replacement
[ ] friendly custom car URLs, e.g. car.Name-car.Year-car.Engine-IDcar.ID
[ ] Add randomly generated mileage and car price for a single car handler or add to json
[ ] Maybe map webapi cars to have full data, and also cache it. Or create a car review struct and cache it for catalog.