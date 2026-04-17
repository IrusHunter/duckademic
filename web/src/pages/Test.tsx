// import css from "./Test module.css"

import { useQuery } from "@tanstack/react-query";
import axios from "axios"

const fetchPerson = async () => {
  const response = await axios.get(`http://localhost:10000/curriculum/curriculums`);
  return response.data;
};

export default function Test () {
    const {data, error, isLoading, isError} = useQuery({
        queryKey: ['curriculum'],
        queryFn: fetchPerson,
    });
    
    return(
        <div style={{paddingLeft: "250px"}}>
            {isLoading && <p>Loading</p>}
            {isError && <p>An error occurred: {String(error)}</p>}
            {data && <pre>{JSON.stringify(data, null, 2)}</pre>}
        </div>
    )
}