### this is a food recommendation api

#### endpoints introduction: 

* POST /diets/meals:
> parameter: 
eg: { 
    "food_name": "2 butter cake",
    "where_eaten": "Paris",
    "date_eaten": "2024-12-29",
    "time_eaten": "21-04-19",
    "periods": "1hr"
}

record what did you eat in the certain moment in one day and also the nutritions inside your food(eg: Carbonhydrate, Fat, Protein, Carolie)

* GET /orders/healthy/{evening}(enter morning,afternoon,evening)/{7(all the dinner(breakfast,lunch) you have eaten within the 7 days)}

recommend you some healthy food based on what the meals you ate within n days. 

* GET /diets/{evening}/{7}

collect all the dinner you ate within 7 days. 

* GET /shop/healthy/{NewYork TimeSquare}/{evening}/{7}

collect all the shop selling healthy food nearby the place you provide(here is the NewYork TimeSquare). According to what you eat within 7 past days, recommend the healthy dinner(here is the dinner you can change to breakfast or lunch) for you.  

* POST /public/register: 

provide the user name and password to let you be authenticated to use this system:

eg: { 
    "username": "wang",
    "password": "123456"
}

* POST /public/login: 

provide the user and password to this server and it will send the jwt token back to the client. Then client can utilize it as a authentication key in header and send the request again to server, so the server can verify it in a short period of time instead of checking the username and password in database. 