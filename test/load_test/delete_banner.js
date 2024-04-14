import { fail } from 'k6';
import http from 'k6/http';


const host = "http://localhost:8080";

let nTargets = 30
let segment = 1000


let cnts = [];


for (let i = 0; i < nTargets; i++) {
    cnts.push(i * segment + 1);
}


export let options = {
    stages: [
        { duration: '2s', target: nTargets },
    ],
};


export default function () {

    // 100 == (select count(*) from banners) / nTargets
    let response = http.del(`${host}/banner/${(__VU-1) * 100 + __ITER}`, null, { headers: { 'token': 'Admin' } });
    
    if (response.status !== 204) {

        console.error(`Error: ${response.body} | Status code: ${response.status}`);
    }

}
