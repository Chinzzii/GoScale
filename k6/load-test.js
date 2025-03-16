import { check, sleep } from 'k6';
import http from 'k6/http';

export let options = {
  stages: [
    { duration: '2m', target: 500 },
    { duration: '5m', target: 1000 },
    { duration: '2m', target: 0 },
  ],
};

export default function () {
  let resA = http.get('http://localhost:30007/a');
  check(resA, {
    'status 200': (r) => r.status === 200,
  });

  sleep(1);
}
