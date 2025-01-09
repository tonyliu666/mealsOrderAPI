## Food Recommendation API

This API provides functionalities for managing and recommending food based on user dietary habits.

**Endpoints:**

* **POST /diets/meals**

    * **Parameters:**
        * `food_name`: Name of the food.
        * `where_eaten`: Location where the food was eaten.
        * `date_eaten`: Date when the food was eaten (YYYY-MM-DD format).
        * `time_eaten`: Time when the food was eaten (HH-MM-SS format).
        * `periods`: Duration of the meal.
        * `nutritions`: Nutritional information (e.g., Carbohydrate, Fat, Protein, Calories).

    * **Description:**
        Records a meal with its details and nutritional information.

* **GET /orders/healthy/{evening}/{7}`**

    * **Parameters:**
        * `{evening}`: Meal time (e.g., "morning", "afternoon", "evening").
        * `{7}`: Number of days to consider (e.g., "7" for the last 7 days).

    * **Description:**
        Recommends healthy food options based on the meals eaten within the specified number of days.

* **GET /diets/{evening}/{7}`**

    * **Parameters:**
        * `{evening}`: Meal time (e.g., "morning", "afternoon", "evening").
        * `{7}`: Number of days to consider (e.g., "7" for the last 7 days).

    * **Description:**
        Retrieves all meals eaten within the specified number of days for the given meal time.

* **GET /shop/healthy/{NewYork TimeSquare}/{evening}/{7}`**

    * **Parameters:**
        * `{NewYork TimeSquare}`: Location where to find healthy food options.
        * `{evening}`: Meal time (e.g., "morning", "afternoon", "evening").
        * `{7}`: Number of days to consider (e.g., "7" for the last 7 days).

    * **Description:**
        Finds shops selling healthy food near the specified location and recommends healthy options for the given meal time based on the user's eating habits within the specified number of days.

* **POST /public/register**

    * **Parameters:**
        * `username`: User's username.
        * `password`: User's password.

    * **Description:**
        Registers a new user with the provided username and password.

* **POST /public/login**

    * **Parameters:**
        * `username`: User's username.
        * `password`: User's password.

    * **Description:**
        Authenticates the user and returns a JWT token for subsequent API requests.


Usage: 


* create .env file: 

like the field defined in the example.env file: 

PORT: (default 8080)
APPID: you should apply this id from this website first: https://developer.edamam.com/edamam-nutrition-api-demo  

APPKeys: the same as above

GEMINI_API_KEY: you can follow this tutorial: https://ai.google.dev/gemini-api/docs/api-key

GOOGLE_MAPS_API_KEY: Here is the tutorial for you to create this key: https://developers.google.com/maps/documentation/javascript/get-api-key 

DBuser,DBpassword, DBport, DBname, DBmode: you can arbitrary specify them 

Then Run it:

> go run main.go 

