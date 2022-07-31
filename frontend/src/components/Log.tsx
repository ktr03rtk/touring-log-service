import { GoogleMap, LoadScript, Marker, Polyline } from '@react-google-maps/api';
import { useState, useEffect } from 'react';
import DatePicker from 'react-datepicker';
import { useNavigate } from 'react-router-dom';

import 'react-datepicker/dist/react-datepicker.css';
import 'react-confirm-alert/src/react-confirm-alert.css';

type LogProperties = {
  jwt: string;
};

type Trip = {
  lat: number;
  lng: number;
};

type Photo = {
  id: string;
  lat: number;
  lng: number;
};

type TouringLog = {
  trip: Trip;
  photo: Photo;
};

const Log = ({ jwt }: LogProperties) => {
  const navigate = useNavigate();
  const [year, setYear] = useState(new Date().getFullYear());
  const [month, setMonth] = useState(new Date().getMonth());
  const [includeDate, setIncludeDate] = useState<Date[]>([]);
  const [photoMarker, setPhotoMarker] = useState<Photo[]>([]);
  const [startDate, setStartDate] = useState(new Date());

  useEffect(() => {
    const payload = `
  {
    dateList(year: ${year}, month: ${month + 1}) {
      day
    }
  }
  `;

    const myHeaders = new Headers();
    myHeaders.append('Content-Type', 'application/json');
    myHeaders.append('Authorization', 'Bearer ' + jwt);

    const requestOptions = {
      method: 'POST',
      body: payload,
      headers: myHeaders,
    };

    fetch(`${process.env.REACT_APP_API_URL}/v1/graphql`, requestOptions)
      .then((response) => response.json())
      .then((data) => {
        return Object.values(data.data.dateList);
      })
      .then((days) => {
        const dates = days.map((m) => new Date(year, month, (m as { day: number }).day));
        dates.push(new Date(year, month - 1, 1));
        dates.push(new Date(year, month + 1, 20));
        setIncludeDate(dates);
        return;
      });
  }, [year, month]);

  useEffect(() => {
    const payload = `
  {
    touringLog(year: ${year}, month: ${month + 1}, day: ${startDate.getDate()}) {
      photo {
        id
        lat
        lng
      }
    }
  }
  `;

    const myHeaders = new Headers();
    myHeaders.append('Content-Type', 'application/json');
    myHeaders.append('Authorization', 'Bearer ' + jwt);

    const requestOptions = {
      method: 'POST',
      body: payload,
      headers: myHeaders,
    };

    fetch(`${process.env.REACT_APP_API_URL}/v1/graphql`, requestOptions)
      .then((response) => response.json())
      .then((data) => {
        const photo = Object.values(data.data.touringLog.photo);
        setPhotoMarker(photo as Photo[]);
        return;
      });
  }, [startDate]);

  const handleChange = (e: any) => {
    setStartDate(e);
  };

  const containerStyle = {
    width: '100%',
    height: '100vh',
  };

  const center = {
    lat: 35.69575,
    lng: 139.77521,
  };

  const onLoad = (polyline: any) => {
    console.log('polyline: ', polyline);
  };

  const path = [
    { lat: 37.772, lng: -122.214 },
    { lat: 21.291, lng: -157.821 },
    { lat: -18.142, lng: 178.431 },
    { lat: -27.467, lng: 153.027 },
  ];

  const options = {
    strokeColor: '#FF0000',
    strokeOpacity: 0.8,
    strokeWeight: 2,
    fillColor: '#FF0000',
    fillOpacity: 0.35,
    clickable: false,
    draggable: false,
    editable: false,
    visible: true,
    radius: 30000,
    paths: [
      { lat: 37.772, lng: -122.214 },
      { lat: 21.291, lng: -157.821 },
      { lat: -18.142, lng: 178.431 },
      { lat: -27.467, lng: 153.027 },
    ],
    zIndex: 1,
  };

  useEffect(() => {
    if (jwt === '') {
      navigate('/login');
      return;
    }
  }, [jwt]);

  return (
    <div>
      <h2>Touring log: Trips and Photos</h2>
      <div className='container'>
        <div> {year} </div>
        <div> {month} </div>
      </div>
      <div className='my-3'>
        <DatePicker
          dateFormat='yyyy/MM/dd'
          selected={startDate}
          onYearChange={(date: Date) => setYear(date.getFullYear())}
          onMonthChange={(date: Date) => setMonth(date.getMonth())}
          onChange={handleChange}
          includeDates={includeDate}
        />
      </div>
      <div className='card'>
        <LoadScript googleMapsApiKey={process.env.REACT_APP_MAP_API_KEY ?? ''}>
          <GoogleMap mapContainerStyle={containerStyle} center={center} zoom={16}>
            {photoMarker.map((m) => {
              return <Marker key={m.id} position={m} />;
            })}
            <Polyline onLoad={onLoad} path={path} options={options} />
          </GoogleMap>
        </LoadScript>
      </div>
    </div>
  );
};

export default Log;
