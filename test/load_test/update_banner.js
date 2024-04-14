import http from 'k6/http';

const host = "http://localhost:8080";

export default function (data) {
    let requestBody = JSON.stringify({
        tag_ids: [4, 5, 6],
        feature_id: 8,
        content: "{\"title\": \"updated banner1\"}",
        is_active: false,
    });
    http.patch(`${host}/banner/${data.bannerId}`, requestBody, { headers: { 'token': 'Admin' } });
}
