import { useState, useEffect } from 'react';
import { confirmAlert } from 'react-confirm-alert';
import { useNavigate } from 'react-router-dom';

import 'react-confirm-alert/src/react-confirm-alert.css';
import Alert from './Alert';

import type { AlertType } from '../types/Alert';
import type { ReactConfirmAlertProps } from 'react-confirm-alert';

type UploadProperties = {
  jwt: string;
};

const Upload = ({ jwt }: UploadProperties) => {
  const [images, setImages] = useState<File[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isSelected, setIsSelected] = useState<boolean>(false);
  const [alert, setAlert] = useState<AlertType>({ type: 'd-none', message: '' });
  const navigate = useNavigate();

  useEffect(() => {
    if (jwt === '') {
      navigate('/login');
      return;
    }
  }, [jwt]);

  const confirmOnSubmit = (e: React.SyntheticEvent) => {
    e.preventDefault();

    // INFO: avoid type error. Types of property 'buttons' are incompatible.
    type UnionProp<T, KEY extends keyof T, TYPE> = { [P in keyof T]: P extends KEY ? T[P] | TYPE : T[P] };
    type CustomConfirmAlertProps = UnionProp<ReactConfirmAlertProps, 'buttons', any>;

    const props: CustomConfirmAlertProps = {
      title: 'Upload phots?',
      message: 'Are you sure?',
      buttons: [
        {
          label: 'Yes',
          onClick: () => {
            setIsLoading(true);
            const myHeaders = new Headers();
            myHeaders.append('Authorization', 'Bearer ' + jwt);

            const data = new FormData();

            images.map((image) => {
              data.append('images', image);
            });

            const requestOptions = {
              method: 'POST',
              body: data,
              headers: myHeaders,
            };

            fetch(`${process.env.REACT_APP_API_URL}/v1/photos`, requestOptions)
              .then((res) => res.json())
              .then((data) => {
                if (data.error) {
                  setAlert({ type: 'alert-danger', message: data.error.message });
                  setIsLoading(false);
                } else {
                  setIsLoading(false);
                }
              })
              .catch((err) => {
                setIsLoading(false);
                console.log(err);
              });
          },
        },
        {
          label: 'No',
        },
      ],
    };

    confirmAlert(props);
  };

  const handleOnAddImage = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files) {
      setIsSelected(false);
      return;
    }
    setImages([...images, ...e.target.files]);
    setIsSelected(true);
  };

  return (
    <>
      <div className='text-center'>
        <h2>Upload your photos.</h2>
        <Alert alertType={alert.type} alertMessage={alert.message} />

        <br />

        <form onSubmit={confirmOnSubmit}>
          <label htmlFor='id1' className='d-grid gap-2 col-6 mx-auto'>
            <span className='btn btn-primary'>SELECT</span>
            <input
              style={{ display: 'none' }}
              id='id1'
              type='file'
              multiple
              accept='image/*,.png,.jpg,.jpeg,.gif'
              onChange={handleOnAddImage}
            />
          </label>

          <br />

          <div className='container'>
            <div className='row'>
              {images.map((image, i) => (
                <img
                  key={i}
                  className='rounded mx-auto d-block col-3'
                  style={{ width: '30%', height: 'auto' }}
                  src={URL.createObjectURL(image)}
                />
              ))}
            </div>
          </div>

          <br />

          {isLoading ? (
            <button className='btn btn-primary  gap-2 col-6 mx-auto' type='submit' disabled>
              <span className='spinner-border spinner-border-sm' role='status' aria-hidden='true'></span>
              UPLOADING...
            </button>
          ) : (
            <button className='btn btn-primary  d-grid gap-2 col-6 mx-auto' type='submit' disabled={!isSelected}>
              UPLOAD
            </button>
          )}
        </form>
      </div>
    </>
  );
};

export default Upload;
