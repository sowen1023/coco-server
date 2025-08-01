import { Avatar, Button, Form, Input, Switch, Upload } from 'antd';
import '../index.scss';
import "./ChatStartPage.scss"
import AIAssistantSelect from '@/pages/ai-assistant/modules/AIAssistantSelect';
import { PlusOutlined } from '@ant-design/icons';

const ChatStartPage = memo((props) => {
  const { startPageSettings, logo, setLogo, isSub, assistants } = props;
  const { t } = useTranslation();
  const { defaultRequiredRule } = useFormRules();
  const [enabled, setEnabled] = useState(false);
  const [excludedAssistants, setExcludedAssistants] = useState([])

  const calculateCharacterLength = (str) => {
    let length = 0;
    for (let i = 0; i < str.length; i++) {
      const code = str.charCodeAt(i);
      if ((code >= 0x0000 && code <= 0x007F)) {
        length += 1;
      } else {
        length += 2;
      }
    }
    return length;
  };

  const renderIcon = (base64) => {
    if (base64) {
      return (
        <div className="chart-start-page-image css-var-r0 ant-btn">
          <Avatar shape="square" src={base64} />
        </div>
      );
    }
    return null;
  }

  const uploadProps = {
    name: "file",
    action: "",
    accept: "image/*,.svg",
  };

  useEffect(() => {
    setEnabled(startPageSettings?.enabled)
  }, [startPageSettings?.enabled])

  useEffect(() => {
    setExcludedAssistants(startPageSettings?.display_assistants || [])
  }, [startPageSettings?.display_assistants])

  return (
    <>
      {
        isSub ? (
          <Form.Item
            label=" "
            labelCol={isSub ? { span: 0 } : {}} 
            className={enabled ? "" : "mb-0px"}
            help={t('page.settings.app_settings.chat_settings.labels.start_page_placeholder')}
          >
            <Form.Item
              label={t('page.settings.app_settings.chat_settings.labels.start_page')}
              name={['start_page', 'enabled']}
              className="mb-0px"
            >
              <Switch size="small" onChange={setEnabled}/>
            </Form.Item>
          </Form.Item>
        ) : (
          <Form.Item
            label={t('page.settings.app_settings.chat_settings.labels.start_page')}
            name={['start_page', 'enabled']}
            help={(
              <div className="mb-24px">{t('page.settings.app_settings.chat_settings.labels.start_page_placeholder')}</div>
            )}
          >
            <Switch size="small" onChange={setEnabled}/>
          </Form.Item>
        )
      }
      {
        enabled && (
          <>
            <Form.Item
              label=" "
              labelCol={isSub ? { span: 0 } : {}} 
              help={(
                <div className="mb-24px">
                  <div>{t('page.settings.app_settings.chat_settings.labels.logo_placeholder')}</div>
                  <div>{t('page.settings.app_settings.chat_settings.labels.logo_size_placeholder')}</div>
                </div>
              )}
            >
              <div className="mb-8px">
                {t('page.settings.app_settings.chat_settings.labels.logo')}
              </div>
            </Form.Item>
            <Form.Item label=" " labelCol={isSub ? { span: 0 } : {}} name={['logo', 'light']}>
              <div className="mb-8px settings-form-help">
                {t('page.settings.app_settings.chat_settings.labels.logo_light')}
              </div>
              <div style={{ display: "flex", gap: 22 }}>
                {renderIcon(logo.light)}
                <Upload
                  {...uploadProps}
                  showUploadList={false}
                  fileList={logo.lightList}
                  beforeUpload={(file) => {
                    setLogo((state) => ({
                      ...state,
                      lightList: [file],
                      lightLoading: true,
                    }))
                    const reader = new FileReader();
                    reader.readAsDataURL(file);
                    reader.onload = () => {
                      setLogo((state) => ({
                        ...state,
                        lightLoading: false,
                        light: reader.result
                      }))
                    };
                    return false
                  }}
                >
                  <Button loading={logo.lightLoading} icon={<SvgIcon className="text-12px" icon="mdi:upload" />}>{t('common.upload')}</Button>
                </Upload>
                <Button className="px-0" type="link" onClick={() => {
                  setLogo((state) => ({
                    ...state,
                    lightLoading: false,
                    light: undefined
                  }));
                }}>{t('common.reset')}</Button>
              </div>
            </Form.Item>
            <Form.Item label=" " labelCol={isSub ? { span: 0 } : {}} name={['logo', 'dark']}>
              <div className="mb-8px settings-form-help">
                {t('page.settings.app_settings.chat_settings.labels.logo_dark')}
              </div>
              <div style={{ display: "flex", gap: 22 }}>
                {renderIcon(logo.dark)}
                <Upload
                  {...uploadProps}
                  showUploadList={false}
                  fileList={logo.darkList}
                  beforeUpload={(file) => {
                    setLogo((state) => ({
                      ...state,
                      darkList: [file],
                      darkLoading: true,
                    }))
                    const reader = new FileReader();
                    reader.readAsDataURL(file);
                    reader.onload = () => {
                      setLogo((state) => ({
                        ...state,
                        darkLoading: false,
                        dark: reader.result
                      }))
                    };
                    return false
                  }}
                >
                  <Button loading={logo.darkLoading} icon={<SvgIcon className="text-12px" icon="mdi:upload" />}>{t('common.upload')}</Button>
                </Upload>
                <Button className="px-0" type="link" onClick={() => {
                  setLogo((state) => ({
                    ...state,
                    darkLoading: false,
                    dark: undefined
                  }))
                }}>{t('common.reset')}</Button>
              </div>
            </Form.Item>
            <Form.Item
              label=" "
              labelCol={isSub ? { span: 0 } : {}} 
            >
              <div className="mb-8px">
                {t('page.settings.app_settings.chat_settings.labels.introduction')}
              </div>
              <Form.Item 
                label=" "
                name={['start_page', 'introduction']}
                help={t('page.settings.app_settings.chat_settings.labels.introduction_placeholder')}
                labelCol={{ span: 0 }}
                rules={[
                  () => ({
                    validator(_, value) {
                      if (!value || calculateCharacterLength(value) <= 60) {
                        return Promise.resolve();
                      }
                      return Promise.reject();
                    },
                  })
                ]}
              >
                <Input.TextArea rows={3} />
              </Form.Item>
            </Form.Item>
            <Form.Item
              label=" "
              labelCol={isSub ? { span: 0 } : {}} 
              className={isSub ? "mb-0px" : ""} 
            >
              <div className="mb-8px">
                {t('page.settings.app_settings.chat_settings.labels.assistant')}
              </div>
              <Form.Item className="mb-0px">
                <Form.List name={['start_page', 'display_assistants']}>
                  {(fields, { add, remove }) => (
                    <>
                      {fields.map((field, index) => {
                        const { key, ...restField } = field;
                        return (
                          <Form.Item key={index} className="mb-0px">
                            <div className="flex gap-6px">
                              <Form.Item
                                {...restField}
                                rules={[defaultRequiredRule]}
                                className="flex-1 mb-8px"
                              >
                                <AIAssistantSelect 
                                  assistants={assistants} 
                                  excluded={excludedAssistants}
                                  onChange={(item) => {
                                    const newExcluded = excludedAssistants.concat([])
                                    newExcluded[index] = item?.id
                                    setExcludedAssistants(newExcluded)
                                  }}
                                />
                              </Form.Item>
                              <Form.Item className="mb-8px">
                                <span onClick={() => {
                                  remove(field.name)
                                  const newExcluded = excludedAssistants.concat([])
                                  newExcluded.splice(field.name)
                                  setExcludedAssistants(newExcluded)
                                }}><SvgIcon className="text-16px cursor-pointer" icon="mdi:minus-circle-outline" /></span>
                              </Form.Item>
                            </div>
                          </Form.Item>
                        )
                      })}
                      <Form.Item className="mb-0px">
                        <Button className="!w-80px" type="primary" disabled={fields.length >= 8} icon={<PlusOutlined />} onClick={() => add()}></Button>
                      </Form.Item>
                    </>
                  )}
                </Form.List>
              </Form.Item>
            </Form.Item>
          </>
        )
      }
      
    </>
  );
});

export default ChatStartPage;