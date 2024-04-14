import { fail } from 'k6';
import http from 'k6/http';

const host = "http://localhost:8080";

let nTargets = 30
let segment = 1000


let cnts = [];
let idx = 1;


for (let i = 0; i < nTargets; i++) {
    cnts.push(i * segment + 1);
}



export let options = {
    stages: [
        { duration: '7s', target: nTargets },
    ],
};

export default function () {
    
    let cnt = cnts[__VU - 1]++
    if (cnt > __VU  * segment) {
        fail("out of segment!!! make segmet bigger\n\n")
    }
    idx++

    let requestBody = JSON.stringify({
        tag_ids: [cnt + 1, cnt + 2, cnt + 3],
        feature_id: cnt+10,
        content: "{\"title\": \"banner1\"}",
        is_active: true,
    });


    let response = http.patch(`${host}/banner/${idx}`, requestBody, { headers: { 'token': 'Admin' } });


    if (response.status !== 200) {
        console.error(`Error: ${response.body} | Status code: ${response.status}`);
    }
}
