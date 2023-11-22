import grpc from 'k6/net/grpc';
import { sleep } from 'k6';

const client = new grpc.Client();
client.load(['definitions'], '../../../../api/api.proto');

export default function () {

  if (__ITER == 0) {
    client.connect('localhost:8081', {
      plaintext: true,
      reflect: false
    });
  }

  // Create a restaurant.
  const restaurant = '{"restaurant":{"name":"Los Baltazares","address":{"street":"Avenida de los pinos de Montequinto","city":"Dos Hermanas","state":"Sevilla","zip":"41089"},"menu":{"items":[{"id":1,"name":"Tapa de adobo","price":"13.14"},{"id":2,"name":"Tapa de patatas","price":"3.14"}]}}}'
  const response = client.invoke('RestaurantService/CreateRestaurant', JSON.parse(restaurant))
  const restaurantId = response.message.restaurantId;

  sleep(Math.random() * 0.5);

  // Update the menu.
  let randomNumber = Math.random();
  if (randomNumber > 0.5) {
    const menu = '{"items":[{"id":1,"name":"Tapa de gambas","price":"3.95"},{"id":2,"name":"Tapa de calamares","price":"6.99"},{"id":3,"name":"Tapa de acelgas","price":"2.99"}]}';
    client.invoke('RestaurantService/UpdateMenu', {restaurant_id: restaurantId, menu: JSON.parse(menu)})
  }

  sleep(Math.random() * 0.5);

  // Delete the restaurant.
  randomNumber = Math.random();
  if (randomNumber > 0.01 && randomNumber < 0.3) {
    client.invoke('RestaurantService/DeleteRestaurant', {restaurant_id: restaurantId})
  }

  //client.close();
}

