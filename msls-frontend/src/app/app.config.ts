import { ApplicationConfig, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';

import { routes } from './app.routes';
import {
  authInterceptor,
  tenantInterceptor,
  requestIdInterceptor,
  errorInterceptor,
  loadingInterceptor,
} from './core/interceptors';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideRouter(routes),
    provideHttpClient(
      withInterceptors([
        // Order matters: request-id first, then tenant, then auth
        // Error and loading wrap all requests
        requestIdInterceptor,
        tenantInterceptor,
        authInterceptor,
        errorInterceptor,
        loadingInterceptor,
      ])
    ),
  ],
};
