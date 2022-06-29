import { Route, Routes, Link } from 'react-router-dom';

import './App.css';
import Home from './components/Home';
import Nav from './components/Nav';
import Upload from './components/Upload';

function App() {
  return (
    <>
      <div className='container'>
        <div className='row'>
          <Nav />
          <div className='col-md-10'>
            <Routes>
              <Route path='/' element={<Home />} />
              <Route path='/upload' element={<Upload />} />
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
