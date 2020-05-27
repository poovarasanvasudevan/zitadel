import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { ProjectRole, User } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

export enum CreationType {
    PROJECT_OWNED = 0,
    PROJECT_GRANTED = 1,
    ORG = 2,
}
@Component({
    selector: 'app-project-member-create-dialog',
    templateUrl: './project-member-create-dialog.component.html',
    styleUrls: ['./project-member-create-dialog.component.scss'],
})
export class ProjectMemberCreateDialogComponent {
    public projectId: string = '';
    public creationType!: CreationType;
    public users: Array<User.AsObject> = [];
    public roles: Array<ProjectRole.AsObject> | string[] = [];
    public CreationType: any = CreationType;
    public memberRoleOptions: string[] = [];
    constructor(
        private projectService: ProjectService,
        public dialogRef: MatDialogRef<ProjectMemberCreateDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
        toastService: ToastService,
    ) {
        this.creationType = data.creationType;
        this.projectId = data.projectId;

        if (this.creationType === CreationType.PROJECT_GRANTED) {
            this.projectService.GetProjectGrantMemberRoles().then(resp => {
                this.memberRoleOptions = resp.toObject().rolesList;
            }).catch(error => {
                toastService.showError(error.message);
            });
        } else if (this.creationType === CreationType.PROJECT_OWNED) {
            this.projectService.GetProjectMemberRoles().then(resp => {
                this.memberRoleOptions = resp.toObject().rolesList;
                console.log(this.memberRoleOptions);
            }).catch(error => {
                toastService.showError(error.message);
            });
        }
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close({ users: this.users, roles: this.roles });
    }

    public setOrgMemberRoles(roles: string[]): void {
        this.roles = roles;
    }
}