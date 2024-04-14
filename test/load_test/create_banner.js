

import { fail } from 'k6';
import http from 'k6/http';

const host = "http://localhost:8080";

let nTargets = 50
let segment = 1000


let cnts = [];


for (let i = 0; i < nTargets; i++) {
    cnts.push(i * segment + 1);
}




export let options = {
    stages: [
        { duration: '3s', target: nTargets },
    ],
};

export default function () {
    
    let cnt = cnts[__VU - 1]++
    if (cnt > __VU  * segment) {
        fail("out of segment!!! make segmet bigger\n\n")
    }

    let requestBody = JSON.stringify({
        tag_ids: [cnt, cnt + 1, cnt + 2],
        feature_id: cnt,
        content: "{\"title\": \"banner1\"}",
        is_active: true,
    });


    let response = http.post(`${host}/banner`, requestBody, { headers: { 'token': 'Admin' } });


    if (response.status !== 201) {
        console.error(`Error: ${response.body} | Status code: ${response.status}`);
    }
}
