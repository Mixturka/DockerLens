import { useState, useEffect } from 'react'

const BACKEND_API_URL = process.env.REACT_APP_BACKEND_API_URL
const GET_PINGS_INTERVAL_MS = process.env.REACT_APP_PING_FETCH_INTERVAL_MS

if (!BACKEND_API_URL) {
    throw new Error("BACKEND_API_URL is not set! Set up environment please!")
}

export default function PingTable() {
    const [pings, setPings] = useState([])
    const [error, setError] = useState(null);

    useEffect(() => {
        const getPings = async () => { 
            try {
                const url = `${BACKEND_API_URL.replace(/\/$/, '')}/pings`;
                const response = await fetch(url);
                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }
                const data = await response.json();
                setPings(data.pings);
            } catch (error) {
                setError("An error occurred trying to get pings. Please try again later.");
                console.error("Failed to fetch pings: ", error);
            }
        };
        getPings();
        const getPingsInterval = setInterval(getPings, GET_PINGS_INTERVAL_MS);
        return () => clearInterval(getPingsInterval);
    }, []);

    const sortedPings = pings.sort((a, b) => {
        return a.IP.localeCompare(b.IP);
    });


    if (error) {
        return (
            <div className="alert alert-danger" role="alert">
                {error}
            </div>
        );
    }

    return (
        <div className="container">
            <h3 className="text-center text-dark">
                Ping Table
            </h3>
            <table className="table table-bordered table-striped table-hover">
                <thead className="table-dark">
                    <tr>
                        <th className="text-center">IP Address</th>
                        <th className="text-center">Ping Duration(ms)</th>
                        <th className="text-center">Last Successful Ping Date</th>
                        <th className="text-center">Status</th>
                    </tr>
                </thead>
                <tbody className="text-center">
                    {sortedPings.map((ping, index) => (
                        <tr key={index}>
                            <td>{ping.IP}</td>
                            <td>{ping.Duration}</td>
                            <td>{new Date(ping.LastSuccess).toLocaleString()}</td>
                            <td class={ping.IsSuccess ? "table-success" : "table-danger"}>
                                {ping.IsSuccess ? `Success` : `Failure`}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};