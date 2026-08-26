import { ErrorBoundary, render } from 'solid-js/web'

import { EnsureAuthenticateRoute } from '@features/auth/ui/EnsureAuthenticateRoute'
import { GroupView } from '@pages/GroupView'
import { GroupsView } from '@pages/GroupsView'
import { HomeView } from '@pages/HomeView'
import { LoginView } from '@pages/LoginView'
import { NotFoundView } from '@pages/NotFoundView'
import { UsrsView } from '@pages/UsrsView'
import { Navigate, Route, Router } from '@solidjs/router'

import { AuthProvider, useAuth } from './features/auth/providers/AuthProvider'
import './index.css'
import { NodeView } from './pages/NodeView'
import { ErrView, Navbar, Toaster } from './shared/ui'

function ProtectedLayout(props) {
    const { state, logout } = useAuth()

    return (
        <section class="flex flex-col h-screen">
            <Navbar
                routes={Array.from(state.routes.values())}
                username={state.usr.username}
                logout={logout}
            />
            <main class="flex-1 overflow-y-auto bg-bg">{props.children}</main>
        </section>
    )
}

export function App() {
    return (
        <AuthProvider>
            <ErrorBoundary fallback={error => <ErrView error={error} />}>
                <Router>
                    <Route path="/login" component={LoginView} />
                    <Route
                        path="/"
                        component={() => <Navigate href="/nodes" />}
                    />
                    <Route path="/nodes">
                        <Route
                            path="/"
                            component={() => (
                                <EnsureAuthenticateRoute>
                                    <ProtectedLayout>
                                        <HomeView />
                                    </ProtectedLayout>
                                </EnsureAuthenticateRoute>
                            )}
                        />
                        <Route
                            path="/:id"
                            component={() => (
                                <EnsureAuthenticateRoute>
                                    <ProtectedLayout>
                                        <NodeView />
                                    </ProtectedLayout>
                                </EnsureAuthenticateRoute>
                            )}
                        />
                    </Route>
                    <Route path="/groups">
                        <Route
                            path="/"
                            component={() => (
                                <EnsureAuthenticateRoute>
                                    <ProtectedLayout>
                                        <GroupsView />
                                    </ProtectedLayout>
                                </EnsureAuthenticateRoute>
                            )}
                        />
                        <Route
                            path="/:id"
                            component={() => (
                                <EnsureAuthenticateRoute>
                                    <ProtectedLayout>
                                        <GroupView />
                                    </ProtectedLayout>
                                </EnsureAuthenticateRoute>
                            )}
                        />
                    </Route>
                    <Route
                        path="/usrs"
                        component={() => (
                            <EnsureAuthenticateRoute>
                                <ProtectedLayout>
                                    <UsrsView />
                                </ProtectedLayout>
                            </EnsureAuthenticateRoute>
                        )}
                    />
                    <Route
                        path="*"
                        component={() => (
                            <EnsureAuthenticateRoute>
                                <ProtectedLayout>
                                    <NotFoundView />
                                </ProtectedLayout>
                            </EnsureAuthenticateRoute>
                        )}
                    />
                </Router>
                <Toaster />
            </ErrorBoundary>
        </AuthProvider>
    )
}

render(App, document.getElementById('root'))
