import { useState } from 'react';
import { confirmAlert } from 'react-confirm-alert';

import 'react-confirm-alert/src/react-confirm-alert.css';
import Alert from './Alert';

import type { AlertType } from '../types/Alert';
import type { ReactConfirmAlertProps } from 'react-confirm-alert';

const Upload = () => {
  const [images, setImages] = useState<File[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [alert, setAlert] = useState<AlertType>({ type: 'd-none', message: '' });

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

            const data = new FormData();

            images.map((image) => {
              data.append('images', image);
            });

            const requestOptions = {
              method: 'POST',
              body: data,
            };

            fetch('http://192.168.10.104:8080/v1/upload', requestOptions)
              .then((res) => res.json())
              .then((data) => {
                if (data.error) {
                  setAlert({ type: 'alert-danger', message: data.error.message });
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

  const handleOnAddImage = (e: any) => {
    if (!e.target.files) return;
    setImages([...images, ...e.target.files]);
  };

  return (
    <>
      <div className='text-center'>
        <h2>Upload your photos.</h2>
        <Alert alertType={alert.type} alertMessage={alert.message} />

        <br />

        <form onSubmit={confirmOnSubmit}>
          <label htmlFor='id1' className='d-grid gap-2 col-6 mx-auto'>
            <span className='btn btn-primary opacity-75'>SELECT</span>
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
            <button className='btn btn-primary opacity-75 gap-2 col-6 mx-auto' type='submit' disabled>
              <span className='spinner-border spinner-border-sm' role='status' aria-hidden='true'></span>
              UPLOADING...
            </button>
          ) : (
            <button className='btn btn-primary opacity-75 d-grid gap-2 col-6 mx-auto' type='submit'>
              UPLOAD
            </button>
          )}
        </form>
      </div>
    </>
  );
};

export default Upload;
