import { useState, useEffect } from 'react';
import { Route, Routes, Link } from 'react-router-dom';

import './App.css';
import Home from './components/Home';
import Log from './components/Log';
import Login from './components/Login';
import Nav from './components/Nav';
import Signup from './components/Signup';
import Upload from './components/Upload';

function App() {
  const [jwt, setJwt] = useState('');

  const handleJWTChange = (token: string) => {
    setJwt(token);
  };

  useEffect(() => {
    const t = window.localStorage.getItem('jwt');
    if (t) {
      if (jwt === '') {
        setJwt(JSON.parse(t));
      }
    }
  }, []);

  return (
    <>
      <div className='container'>
        <div className='row'>
          <Nav jwt={jwt} setJwt={setJwt} />
          <div className='col-md-10'>
            <Routes>
              <Route path='/' element={<Home />} />
              <Route path='/log' element={<Log jwt={jwt} />} />
              <Route path='/upload' element={<Upload jwt={jwt} />} />
              <Route path='/signup' element={<Signup />} />
              <Route path='/login' element={<Login handleJWTChange={handleJWTChange} />} />
              <Route path='*' element={<NoMatch />} />
            </Routes>
          </div>
        </div>
      </div>
    </>
  );
}

export default App;

function NoMatch() {
  return (
    <div>
      <h2>Nothing to see here!</h2>
      <p>
        <Link to='/'>Go to the home page</Link>
      </p>
    </div>
  );
}
