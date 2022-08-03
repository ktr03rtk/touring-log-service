import { GoogleMap, LoadScript, Marker, Polyline } from '@react-google-maps/api';
import { useState, useEffect } from 'react';
import DatePicker from 'react-datepicker';
import { useNavigate } from 'react-router-dom';

import Alert from './Alert';
import Modal from './Modal';

import type { AlertType } from '../types/Alert';

import 'react-confirm-alert/src/react-confirm-alert.css';
import 'react-datepicker/dist/react-datepicker.css';

type LogProperties = {
  jwt: string;
};

type TripTypes = {
  lat: number;
  lng: number;
};

type PhotoTypes = {
  id: string;
  lat: number;
  lng: number;
};

const Log = ({ jwt }: LogProperties) => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [year, setYear] = useState(new Date().getFullYear());
  const [month, setMonth] = useState(new Date().getMonth());
  const [includeDate, setIncludeDate] = useState<Date[]>([]);
  const [photoMarker, setPhotoMarker] = useState<PhotoTypes[]>([]);
  const [startDate, setStartDate] = useState(new Date());
  const [alert, setAlert] = useState<AlertType>({ type: 'd-none', message: '' });
  const [image, setImage] = useState('');
  const [show, setShow] = useState(false);

  const CustomMarker = (props: any) => {
    const { id } = props;

    const onMarkerClick = () => {
      setShow(true);
      setIsLoading(true);
      const myHeaders = new Headers();
      myHeaders.append('Content-Type', 'application/json');
      myHeaders.append('Authorization', 'Bearer ' + jwt);

      const requestOptions = {
        method: 'GET',
        headers: myHeaders,
      };

      fetch(`${process.env.REACT_APP_API_URL}/v1/photos/` + id, requestOptions)
        .then((res) => res.json())
        .then((data) => {
          if (data.error) {
            setAlert({ type: 'alert-danger', message: data.error.message });
            setIsLoading(false);
            return;
          }
          const img = Object.values(data.response.message);
          return img;
        })
        .then((img) => {
          setImage('data:image/png;base64, ' + (img as string[]).join(''));
          setIsLoading(false);
        })
        .catch((err) => {
          setIsLoading(false);
          console.log(err);
        });
    };
    return <Marker onClick={onMarkerClick} {...props} />;
  };

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
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          setAlert({ type: 'alert-danger', message: data.error.message });
          return;
        }
        const photo = Object.values(data.data.touringLog.photo);
        setPhotoMarker(photo as PhotoTypes[]);
        return;
      })
      .catch((err) => {
        console.log(err);
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
      <Alert alertType={alert.type} alertMessage={alert.message} />

      <br />

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
              return <CustomMarker key={m.id} id={m.id} position={m} />;
            })}
            <Polyline onLoad={onLoad} path={path} options={options} />
          </GoogleMap>
        </LoadScript>
      </div>
      <Modal
        show={show}
        title='Photo'
        body={<>{isLoading ? <p>Downloading...</p> : <img src={image} width='80%' height='80%' />}</>}
        footer={
          <button
            onClick={() => {
              setShow(false);
            }}
          >
            Close
          </button>
        }
        onHide={() => {
          setShow(false);
        }}
      />
    </div>
  );
};

export default Log;
