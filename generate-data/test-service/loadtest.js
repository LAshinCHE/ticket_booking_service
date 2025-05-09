import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 10,           
  duration: '10s'    
};

export default function () {
  const ticket_id = Math.floor(Math.random() * 1000);       
  const user_id = Math.floor(Math.random() * 1000);         
  const payment = Math.floor(Math.random() * 1000) + 100;   

  const booking_id = (__VU - 1) * 100000 + __ITER;

  const payload = JSON.stringify({
    booking_id: booking_id,
    ticket_id: ticket_id,
    user_id: user_id,
    payment: payment
  });

  const headers = { 'Content-Type': 'application/json' };

  http.post('http://localhost:8081/booking/', payload, { headers });

  sleep(1); 
}
