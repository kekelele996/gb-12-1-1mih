import { createBrowserRouter, Navigate } from 'react-router-dom'
import App from '@/App'
import Dashboard from '@/pages/Dashboard'
import UniversityList from '@/pages/UniversityList'
import UniversityDetail from '@/pages/UniversityDetail'
import ApplicationList from '@/pages/ApplicationList'
import ApplicationDetail from '@/pages/ApplicationDetail'
import DocumentEditor from '@/pages/DocumentEditor'
import TimelinePage from '@/pages/TimelinePage'
import Messages from '@/pages/Messages'
import CounselorWorkbench from '@/pages/CounselorWorkbench'
import Profile from '@/pages/Profile'
import Login from '@/pages/Login'
import RoleGuard from '@/components/common/RoleGuard'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: 'dashboard', element: <RoleGuard roles={['student', 'counselor', 'admin']}><Dashboard /></RoleGuard> },
      { path: 'universities', element: <UniversityList /> },
      { path: 'universities/:id', element: <UniversityDetail /> },
      { path: 'applications', element: <RoleGuard roles={['student', 'counselor', 'admin']}><ApplicationList /></RoleGuard> },
      { path: 'applications/:id', element: <RoleGuard roles={['student', 'counselor', 'admin']}><ApplicationDetail /></RoleGuard> },
      { path: 'documents/:id', element: <RoleGuard roles={['student', 'counselor', 'admin']}><DocumentEditor /></RoleGuard> },
      { path: 'timeline', element: <RoleGuard roles={['student', 'counselor']}><TimelinePage /></RoleGuard> },
      { path: 'messages', element: <RoleGuard roles={['student', 'counselor']}><Messages /></RoleGuard> },
      { path: 'counselor', element: <RoleGuard roles={['counselor', 'admin']}><CounselorWorkbench /></RoleGuard> },
      { path: 'profile', element: <RoleGuard roles={['student', 'counselor', 'admin']}><Profile /></RoleGuard> },
      { path: 'login', element: <Login /> },
      { path: '*', element: <Navigate to="/" replace /> },
    ],
  },
])
