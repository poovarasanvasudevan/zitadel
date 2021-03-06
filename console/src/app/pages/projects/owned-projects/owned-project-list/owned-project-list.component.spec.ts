import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GrantedProjectListComponent } from './granted-project-list.component';

describe('ProjectListComponent', () => {
    let component: GrantedProjectListComponent;
    let fixture: ComponentFixture<GrantedProjectListComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [GrantedProjectListComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(GrantedProjectListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
