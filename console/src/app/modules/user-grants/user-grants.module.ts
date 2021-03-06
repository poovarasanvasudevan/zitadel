import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { AvatarModule } from '../avatar/avatar.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { UserGrantsComponent } from './user-grants.component';



@NgModule({
    declarations: [UserGrantsComponent],
    imports: [
        CommonModule,
        FormsModule,
        AvatarModule,
        MatButtonModule,
        HasRoleModule,
        MatTableModule,
        MatPaginatorModule,
        MatIconModule,
        RouterModule,
        MatProgressSpinnerModule,
        MatCheckboxModule,
        MatTooltipModule,
        MatSelectModule,
        MatFormFieldModule,
        TranslateModule,
        HasRolePipeModule,
        TimestampToDatePipeModule,
        RefreshTableModule,
        LocalizedDatePipeModule,
    ],
    exports: [
        UserGrantsComponent,
    ],
})
export class UserGrantsModule { }
