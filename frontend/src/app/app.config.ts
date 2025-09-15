import { ApplicationConfig, provideBrowserGlobalErrorListeners, provideZonelessChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideStore } from '@ngrx/store';
import { provideEffects } from '@ngrx/effects';
import { provideStoreDevtools } from '@ngrx/store-devtools';
import { provideNzI18n, en_US } from 'ng-zorro-antd/i18n';
import { provideNzIcons } from 'ng-zorro-antd/icon';
import { IconDefinition } from '@ant-design/icons-angular';
import { 
  UserOutline, 
  LockOutline, 
  MailOutline, 
  EyeInvisibleOutline, 
  EyeTwoTone,
  DownOutline,
  MenuOutline,
  LogoutOutline,
  UserAddOutline,
  CheckCircleOutline,
  PlayCircleOutline,
  FireOutline,
  ThunderboltOutline,
  TrophyOutline,
  CodeOutline,
  FilterOutline,
  SortAscendingOutline,
  BellOutline,
  HistoryOutline,
  ArrowLeftOutline,
  EyeOutline,
  ReloadOutline,
  CloseCircleOutline,
  ClockCircleOutline,
  MessageOutline
} from '@ant-design/icons-angular/icons';
import { registerLocaleData } from '@angular/common';
import en from '@angular/common/locales/en';

const icons: IconDefinition[] = [
  UserOutline,
  LockOutline,
  MailOutline,
  EyeInvisibleOutline,
  EyeTwoTone,
  DownOutline,
  MenuOutline,
  LogoutOutline,
  UserAddOutline,
  CheckCircleOutline,
  PlayCircleOutline,
  FireOutline,
  ThunderboltOutline,
  TrophyOutline,
  CodeOutline,
  FilterOutline,
  SortAscendingOutline,
  BellOutline,
  HistoryOutline,
  ArrowLeftOutline,
  EyeOutline,
  ReloadOutline,
  CloseCircleOutline,
  ClockCircleOutline,
  MessageOutline
];
import { routes } from './app.routes';
import { authInterceptor } from './auth/interceptors/auth.interceptor';
import { authReducer } from './auth/store/auth.reducer';
import { AuthEffects } from './auth/store/auth.effects';
import { problemReducer } from './problems/store/problem.reducer';
import { ProblemEffects } from './problems/store/problem.effects';
import { NzMessageService } from 'ng-zorro-antd/message';

registerLocaleData(en);

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideZonelessChangeDetection(),
    provideRouter(routes),
    provideHttpClient(withInterceptors([authInterceptor])),
    provideAnimationsAsync(),
    provideNzI18n(en_US),
    provideNzIcons(icons),
    provideStore({ 
      auth: authReducer,
      problems: problemReducer 
    }),
    provideEffects([AuthEffects, ProblemEffects]),
    NzMessageService,
    provideStoreDevtools({ maxAge: 25 })
  ]
};
