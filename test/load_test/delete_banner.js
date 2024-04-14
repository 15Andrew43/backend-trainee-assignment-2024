import http from 'k6/http';

const host = "http://localhost:8080";

export function teardown(data) {
    http.del(`${host}/banner/${data.bannerId}`, null, { headers: { 'token': 'Admin' } });
}
