import http from 'k6/http';

const host = "http://localhost:8080";

export default function () {
    let response = http.get(`${host}/user_banner?tag_id=1&feature_id=5`, { headers: { 'token': 'AuthorizedUser' } });
}
