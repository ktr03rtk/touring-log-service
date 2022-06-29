import { Link } from 'react-router-dom';

const Nav = () => {
  return (
    <nav className='navbar navbar-expand-lg bg-secondary navbar-dark bg-opacity-75'>
      <div className='container-fluid'>
        <Link to='/' className='navbar-brand'>
          Touring log
        </Link>
        <div className='collapse navbar-collapse' id='navbarSupportedContent'>
          <ul className='navbar-nav me-auto mb-2 mb-lg-0'>
            <li className='nav-item'>
              <Link to='/' className='nav-link'>
                Home
              </Link>
            </li>
            <li className='nav-item'>
              <Link to='/upload' className='nav-link'>
                Upload
              </Link>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
};
export default Nav;
