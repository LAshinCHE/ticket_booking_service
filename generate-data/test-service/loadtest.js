import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  scenarios: {
    constant_rate_test: {
      executor: 'constant-arrival-rate',
      rate: 10,               
      timeUnit: '1s',           
      duration: '4s',
      preAllocatedVUs: 5,     
      maxVUs: 5,            
    },
  },
};

export default function () {
  const ticketID = Math.floor(Math.random() * 1000) + __VU + __ITER;
  const userID = Math.floor(Math.random() * 1000) + __VU + __ITER;
  const price = Math.floor(Math.random() * 100) + 40;   


  const payload = JSON.stringify({
    ticketID: ticketID,
    userID: userID,
    price: price
  });

  const headers = { 'Content-Type': 'application/json' };

  http.post('http://localhost:8081/booking/', payload, { headers });
}
