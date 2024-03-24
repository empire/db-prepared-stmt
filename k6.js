import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 20,
  duration: '600s',
};

export default function () {
  const res = http.get('http://localhost:7070/random');
  check(res, {
    'is status 200': (r) => r.status === 200,
    'check body': (r) => r.body && r.body.includes('User'),
  });
}
