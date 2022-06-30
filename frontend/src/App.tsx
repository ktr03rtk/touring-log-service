import { useState } from 'react';
import { Route, Routes, Link } from 'react-router-dom';

import './App.css';
import Home from './components/Home';
import Login from './components/Login';
import Nav from './components/Nav';
import Signup from './components/Signup';
import Upload from './components/Upload';

function App() {
  const [jwt, setJwt] = useState('');

  return (
    <>
      <div className='container'>
        <div className='row'>
          <Nav jwt={jwt} setJwt={setJwt} />
          <div className='col-md-10'>
            <Routes>
              <Route path='/' element={<Home />} />
              <Route path='/upload' element={<Upload />} />
              <Route path='/signup' element={<Signup />} />
              <Route path='/login' element={<Login />} />
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
