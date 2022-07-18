import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import Alert from './Alert';
import Input from './Input';

import type { AlertType } from '../types/Alert';

type FormData = {
  email: string;
  password: string;
  unit: string;
};

const Signup = () => {
  const [formData, setFormData] = useState<FormData>({ email: '', password: '', unit: '' });
  const [errors, setErrors] = useState<string[]>(['']);
  const [alert, setAlert] = useState<AlertType>({ type: 'd-none', message: '' });
  const navigate = useNavigate();

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const err = [];

    if (formData.email === '') {
      err.push('email');
    }

    if (formData.password === '') {
      err.push('password');
    }

    if (formData.password === '') {
      err.push('unit');
    }

    if (err.length > 0) {
      setErrors(err);
      return false;
    }

    const data = new FormData(e.currentTarget);
    const payload = Object.fromEntries(data.entries());
    const myHeaders = new Headers();
    myHeaders.append('Content-Type', 'application/json');

    const requestOptions = {
      method: 'POST',
      body: JSON.stringify(payload),
      headers: myHeaders,
    };

    fetch(`${process.env.REACT_APP_API_URL}/v1/signup`, requestOptions)
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          setAlert({
            type: 'alert-danger',
            message: data.error.message,
          });
        } else {
          navigate('/login');
        }
      });
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const hasError = (key: string) => {
    return errors.indexOf(key) !== -1;
  };

  return (
    <>
      <h2>Signup</h2>
      <hr />
      <Alert alertType={alert.type} alertMessage={alert.message} />

      <form className='pt-3' onSubmit={handleSubmit}>
        <Input
          title={'Email'}
          type={'email'}
          name={'email'}
          handleChange={handleChange}
          className={hasError('email') ? 'is-invalid' : ''}
          errorDiv={hasError('email') ? 'text-danger' : 'd-none'}
          errorMsg={'Please enter a valid email address'}
        />

        <Input
          title={'Password'}
          type={'password'}
          name={'password'}
          handleChange={handleChange}
          className={hasError('password') ? 'is-invalid' : ''}
          errorDiv={hasError('password') ? 'text-danger' : 'd-none'}
          errorMsg={'Please enter a valid password'}
        />

        <Input
          title={'Unit'}
          type={'unit'}
          name={'unit'}
          handleChange={handleChange}
          className={hasError('unit') ? 'is-invalid' : ''}
          errorDiv={hasError('unit') ? 'text-danger' : 'd-none'}
          errorMsg={'Please enter a valid unit'}
        />
        <hr />

        <button className='btn btn-primary'>Signup</button>
      </form>
    </>
  );
};

export default Signup;
