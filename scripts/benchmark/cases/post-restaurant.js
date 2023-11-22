import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  const restaurant = '{"restaurant":{"name":"Los Baltazares","address":{"street":"Avenida de los pinos de Montequinto","city":"Dos Hermanas","state":"Sevilla","zip":"41089"},"menu":{"items":[{"id":1,"name":"Tapa de adobo","price":"13.14"},{"id":2,"name":"Tapa de patatas","price":"3.14"}]}}}'

  // Create a restaurant.
  let res = http.post('http://localhost:8080/api/v1/restaurants', restaurant, {
    headers: { 'Content-Type': 'application/json' },
  });
  const restaurantId = JSON.parse(res.body).restaurantId;
  
  sleep(Math.random() * 0.5);

  // Update the menu.
  let randomNumber = Math.random();
  if (randomNumber > 0.5) {
    const menu = '{"menu":{"items":[{"id":1,"name":"Tapa de gambas","price":"3.95"},{"id":2,"name":"Tapa de calamares","price":"6.99"},{"id":3,"name":"Tapa de acelgas","price":"2.99"}]}}';
    http.put('http://localhost:8080/api/v1/restaurants/' + restaurantId + '/menu', menu, {
      headers: { 'Content-Type': 'application/json' },
    });
  }

  sleep(Math.random() * 0.5);

  // Delete the restaurant.
  randomNumber = Math.random();
  if (randomNumber > 0.01 && randomNumber < 0.3) {
    http.del('http://localhost:8080/api/v1/restaurants/' + restaurantId);
  }
}

