import { fail } from 'k6';
import http from 'k6/http';



const host = "http://localhost:8080";

let nTargets = 50



export let options = {
    stages: [
        { duration: '1s', target: nTargets },
    ],
};


export default function () {


    let response = http.get(`${host}/user_banner?tag_id=${__ITER+1}&feature_id=${__ITER+1}`, { headers: { 'token': 'AuthorizedUser' } });
    
    if (response.status !== 200) {
        console.error(`Error: ${response.body} | Status code: ${response.status}`);
    }

}
