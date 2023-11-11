import http from 'k6/http'

export default function () {
  const data = '{"restaurant":{"name":"Los Baltazares","address":{"street":"Avenida de los pinos de Montequinto","city":"Dos Hermanas","state":"Sevilla","zip":"41089"},"menu":{"items":[{"id":1,"name":"Tapa de adobo","price":"13.14"},{"id":2,"name":"Tapa de patatas","price":"3.14"}]}}}'
  let res = http.post('http://localhost:8080/api/v1/restaurants', data, {
    headers: { 'Content-Type': 'application/json' },
  });
}

