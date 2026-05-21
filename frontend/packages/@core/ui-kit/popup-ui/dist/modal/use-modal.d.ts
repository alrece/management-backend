import { ExtendedModalApi, ModalApiOptions, ModalProps } from "./modal.js";
import * as _$vue from "vue";

//#region src/modal/use-modal.d.ts
declare function setDefaultModalProps(props: Partial<ModalProps>): void;
declare function useVbenModal<TParentModalProps extends ModalProps = ModalProps>(options?: ModalApiOptions): readonly [_$vue.DefineSetupFnComponent<TParentModalProps, {}, {}, TParentModalProps & {}, _$vue.PublicProps>, ExtendedModalApi] | readonly [_$vue.DefineSetupFnComponent<ModalProps, {}, {}, ModalProps & {}, _$vue.PublicProps>, ExtendedModalApi];
//#endregion
export { setDefaultModalProps, useVbenModal };