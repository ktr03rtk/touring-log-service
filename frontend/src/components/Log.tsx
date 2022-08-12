import { useState, useEffect } from 'react';
import DatePicker from 'react-datepicker';
import { useNavigate } from 'react-router-dom';

import Alert from './Alert';
import Map from './Map';

import type { AlertType } from '../types/Alert';
import type { TripTypes, PhotoTypes, TouringLogTypes } from '../types/Touring';

import 'react-confirm-alert/src/react-confirm-alert.css';
import 'react-datepicker/dist/react-datepicker.css';

type LogProperties = {
  jwt: string;
};

const Log = ({ jwt }: LogProperties) => {
  const TokyoStation = { lat: 35.6809155, lng: 139.76606 };
  const navigate = useNavigate();
  const [year, setYear] = useState(new Date().getFullYear());
  const [month, setMonth] = useState(new Date().getMonth());
  const [includeDate, setIncludeDate] = useState<Date[]>([]);
  const [photoMarker, setPhotoMarker] = useState<PhotoTypes[]>([]);
  const [startDate, setStartDate] = useState(new Date());
  const [alert, setAlert] = useState<AlertType>({ type: 'd-none', message: '' });
  const [center, setCenter] = useState<TripTypes>(TokyoStation);
  const [paths, setPaths] = useState<TripTypes[][]>([]);

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
      trip {
        lat
        lng
      }
      photo {
        id
        lat
        lng
      }
      center {
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
        const log = Object.values(data.data)[0] as TouringLogTypes;
        setPhotoMarker(log.photo);
        setPaths([log.trip]);

        if (log.center !== null) {
          setCenter(log.center);
        }
        return;
      })
      .catch((err) => {
        console.log(err);
      });
  }, [startDate]);

  const handleChange = (e: any) => {
    setStartDate(e);
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
      <Map jwt={jwt} photoMarker={photoMarker} setAlert={setAlert} center={center} paths={paths} />
    </div>
  );
};

export default Log;
