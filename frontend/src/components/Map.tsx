import { GoogleMap, LoadScript, Marker, Polyline } from '@react-google-maps/api';
import { useState } from 'react';

import Modal from './Modal';

import type { AlertType } from '../types/Alert';
import type { TripTypes, PhotoTypes } from '../types/Touring';

import 'react-confirm-alert/src/react-confirm-alert.css';
import 'react-datepicker/dist/react-datepicker.css';

type MapProperties = {
  jwt: string;
  photoMarker: PhotoTypes[];
  setAlert: React.Dispatch<React.SetStateAction<AlertType>>;
  center: TripTypes;
  paths: TripTypes[][];
};

const Map = ({ jwt, photoMarker, setAlert, center, paths }: MapProperties) => {
  const [isLoading, setIsLoading] = useState<boolean>(false);
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

  const containerStyle = {
    width: '100%',
    height: '100vh',
  };

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
    paths: paths,
    zIndex: 1,
  };

  return (
    <div>
      <div className='card'>
        <LoadScript googleMapsApiKey={process.env.REACT_APP_MAP_API_KEY ?? ''}>
          <GoogleMap mapContainerStyle={containerStyle} center={center} zoom={12}>
            {paths.map((m) => {
              return <Polyline key={1} path={m} options={options} />;
            })}
            {photoMarker.map((m) => {
              return <CustomMarker key={m.id} id={m.id} position={m} />;
            })}
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

export default Map;
