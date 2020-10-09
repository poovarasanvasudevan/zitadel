import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PasswordComplexityPolicyComponent } from './password-complexity-policy.component';

const routes: Routes = [
    {
        path: '',
        component: PasswordComplexityPolicyComponent,
        data: {
            animation: 'DetailPage',
        },
    },
    // {
    //     path: 'create',
    //     component: PasswordComplexityPolicyComponent,
    //     data: {
    //         animation: 'DetailPage',
    //         action: PolicyComponentAction.CREATE,
    //     },
    // },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class PasswordComplexityPolicyRoutingModule { }
