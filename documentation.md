# Go-Courier API Documentation (gRPC)
 
# AuthService (port:8081)
  ## 1.Register 
  ``Method: POST``<br>
  ``Route: /register`` 
  
  **Request**: `RegisterRequest`
  ```json
    {
      "username": "string",
      "password": "string"
    }
  ```
  
  **Response**: `AuthResponse` 
  ```json
    {
      "message": "string"
    } 
  ```
  
  ## 2. Login   
  ``Method: POST``<br>
  ``Route: /login``
  **Request Type**: `LoginRequest`
  ```json
    {
      "username": "string",
      "password": "string"
    }
  ```
  **Response Type**: `AuthResponse`
  ```json
  {
    "token": "string",
    "message": "string"
  }
  ```


# OrderService (port:8091)
  ## 1.CreateOrder 
  ``Method: POST``<br>
  ``Route: /orders/``

  **Request Type**: `OrderRequest`
  ```json
    {
      "sender_id": "int64",
      "items": "string",
      "address": "string"
    }
  ```
  
  **Response Type**: `OrderResponse`
  ```json
    {
      "id": "int64",
      "status": "string"
    }
  ```

  ## 2.GetOrder
  ``Method: GET``<br>
  ``Route: /orders/{orderId}``

  **Response Type**: `Order`
  ```json
    {
      "id": "int64",
      "sender_id": "int64",
      "items": "string",
      "address": "string",
      "status": "string",
      "created_at": "string"
    }
  ```

  ## 3.ListOrder
  ``Method: GET``<br>
  ``Route: /orders/``

  **Response Type**: `OrderList`
  ```json
    {
      "orders": "Order"
    }
  ```

  ## 4.UpdateOrder
  ``Method: PUT``<br>
  ``Route: /orders/{orderId}``

  **Request Type**: `UpdateOrderRequest`
  ```json
    {
      "id": "int64",
      "items": "string",
      "address": "string"
    }
  ```

  **Response Type**: `UpdateOrderResponse`
  ```json
    {
      "message": "string"
    }
  ```

  ## 5.DeleteOrder
  ``Method: DELETE``<br>
  ``Route: /orders/{orderId}``
  
  **Response Type**: `DeleteOrderResponse`
  ```json
    {
      "message": "string"
    }
  ```


# CourierService (port:8101) ##
  ## 1.RegisterCourier
  ``Method: POST``<br>
  ``Route: /couriers/``

  **Request**: `CourierRequest`
  ```json
    {
      "name": "string",
      "phone_number": "string"
    }
  ```
  
  **Response**: `CourierResponse`
  ```json
    {
      "message": "string"
    }
  ```

  ## 2.GetCourier
  ``Method: GET``<br>
  ``Route: /couriers/{courierId}``
  
  **Response Type**: `Courier`
  ```json
    {
      "id": "int64",
      "name": "string",
      "phone_number": "string",
      "available": "bool",
      "added_by": "int64"
    }
  ```

  ## 3.ListCourier
  ``Method: GET``<br>
  ``Route: /couriers/``
  
  **Response Type**: `CourierList`
  ```json
    {
      "couriers": "Courier"
    }
  ```

  ## 4.UpdateCourier
  ``Method: PUT``<br>
  ``Route: /couriers/{courierId}``
  
  **Request Type**: `UpdateCourierRequest`
  ```json
    {
      "id": "int64",
      "name": "string",
      "phone_number": "string"
    }
  ```
  
  **Response Type**: `CourierResponse`
  ```json
    {
      "message": "string"
    }
  ```

  ## 5.DeleteCourier
  ``Method: DELETE``<br>
  ``Route: /couriers/{courierId}``
  
  **Response Type**: `CourierResponse`
  ```json
    {
      "message": "string"
    }
  ```

  ## 6.ChangeAvailability
  ``Method: PUT``<br>
  ``Route: /couriers/changeAvailability/{courierId}``

  **Request Type**: `ChangeAvailabilityCourierRequest`
  ```json
    {
      "id": "int64"
    }
  ```

  **Response Type**: `CourierResponse`
  ```json
    {
      "message": "string"
    }
  ```

# DeliveryService (port:8111)

  ## 1.AssignCourier
  ``Method: POST``<br>
  ``Route: /assignCourier/``

  **Request Type**: `DeliveryRequest`
  ```json
    {
      "courierId": "int64",
      "orderId": "int64"
    }
  ```

  **Response Type**: `DeliveryResponse`
  ```json
    {
      "message": "string",
      "status": "string"
    }
  ```

  ## 2.CompleteOrder 
  ``Method: PUT``<br>
  ``Route: /completeOrder/{deliveryId}``

  **Response Type**: `DeliveryResponse`
  ```json
    {
      "message": "string",
      "status": "string"
    }
  ```

  ## 3.GetDelivery 
  ``Method: GET``<br>
  ``Route: /deliveries/{deliveryId}``

  **Response Type**: `Delivery`
  ```json
    {
      "id": "int64",
      "courier_data": "Courier",
      "orderData": "Order",
      "status": "string",
      "added_by": "int64"
    }
  ```

  ## 4.ListDelivery 
  ``Method: GET``<br>
  ``Route: /deliveries/``

  **Response Type**: `DeliveryList`
  ```json
    {
      "deliveries": "Delivery"
    }
  ```

  ## 5.DeleteDelivery
  ``Method: DELETE``<br>
  ``Route: /deliveries/{deliveryId}``

  **Response Type**: `DeliveryResponse`
  ```json
    {
      "message": "string",
      "status": "string"
    }
  ```