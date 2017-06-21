
##TODO: 


#Handler
GetNutritionalValue

#Model
AddIngredient
RemoveIngredient
DeleteRecipe
AddRecipe
createSchedule
updateSchedule
AddStep
RemoveStep

#Handler or Model:
RemoveIngredientFromStorage



#### HOUSE METHODS ####

#-------------------------------GET------------------------------#
# GetHouseUsers (house) - DONE
SELECT U.ID, U.EMAIL, U.PASSWORD, U.USERNAME, O.OWN_TYPE, O.DESCRIPTION FROM USER_INFO U INNER JOIN MEMBER_OF M ON M.USER_ID = U.ID
INNER JOIN OWNERSHIP O ON O.OWN_TYPE = M.OWN_TYPE WHERE M.HOUSE_ID = %s
# GetHouseStorage(house)- DONE
SELECT I.NAME, S.AMOUNT, U.NAME FROM INGREDIENT I INNER JOIN ITEM_IN_STORAGE S ON I.ID = S.INGREDIENT_ID 
INNER JOIN UNIT U ON U.ID = S.UNIT_ID WHERE S.HOUSE_ID = %s
#GetHouseRecipes (house) - DONE
SELECT R.ID, R.NAME FROM RECIPE R INNER JOIN HOUSE_RECIPE H ON R.ID = H.RECIPE_ID WHERE H.HOUSE_ID = $1
#GetHouseSchedule(house)- DONE
SELECT W.DAY, T.TYPE, R.NAME FROM RECIPE R, WEEKDAY W, MEAL_TYPE T, SCHEDULE S WHERE S.HOUSE_ID = %s AND S.WEEK_ID = W.ID AND S.TYPE_ID = T.ID AND S.RECIPE_ID = R.ID
#-------------------------------SET------------------------------#
#AddUser


#### USER METHODS ####

#-------------------------------GET------------------------------#
#GetHouses (user) - DONE
SELECT H.ID, H.NAME, O.OWN_TYPE, O.DESCRIPTION FROM HOUSE H INNER JOIN MEMBER_OF M ON M.HOUSE_ID = H.ID 
INNER JOIN OWNERSHIP O ON O.OWN_TYPE = M.OWN_TYPE WHERE M.USER_ID = %s
#GetUserRecipes (user)
SELECT R.ID, R.NAME FROM RECIPE R INNER JOIN USER_RECIPE U ON R.ID = U.RECIPE_ID WHERE U.USER_ID = %s
#GetByUsername(username)
SELECT * FROM USER_INFO I WHERE I.USERNAME = %s
#-------------------------------SET------------------------------#
#AddItemToStorage

#-------------------------------DELETE------------------------------#



#### RECIPE METHODS ####

#-------------------------------GET------------------------------#
#GetSteps (recipe)
SELECT S.ID, S.TEXT FROM STEP S INNER JOIN RECIPE R ON S.RECIPE_ID = R.ID WHERE R.ID = %s
#GetIngredients (recipe)
SELECT * FROM INGREDIENT I INNER JOIN STEP_INGREDIENT S ON S.INGREDIENT_ID = I.ID WHERE S.RECIPE_ID = %s 
#-------------------------------SET------------------------------#



#### INGREDIENT METHODS ####

#-------------------------------GET------------------------------#
# GetIngredientOfStep(recipe, stepNo, ingredient)
SELECT U.NAME, S.AMOUNT FROM UNIT U INNER JOIN STEP_INGREDIENT S ON U.ID = S.UNIT_ID WHERE S.RECIPE_ID = %s AND S.STEP_ID = %s AND S.INGREDIENT_ID = %s 
#-------------------------------SET------------------------------#



